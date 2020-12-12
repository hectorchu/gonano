package pow

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash"
	"math/rand"
	"runtime"

	"github.com/robvanmieghem/go-opencl/cl"
	"golang.org/x/crypto/blake2b"
)

// Generate generates proof-of-work.
func Generate(data, difficulty []byte) (work []byte, err error) {
	target := binary.BigEndian.Uint64(difficulty)
	if work, err = GenerateGPU(data, target); err != nil {
		work, err = GenerateCPU(data, target)
	}
	for i, j := 0, len(work)-1; i < j; i, j = i+1, j-1 {
		work[i], work[j] = work[j], work[i]
	}
	return
}

// GenerateGPU generates proof-of-work using the GPU.
func GenerateGPU(data []byte, target uint64) (work []byte, err error) {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		return
	}
	for _, platform := range platforms {
		devices, err := platform.GetDevices(cl.DeviceTypeGPU)
		if err != nil || len(devices) == 0 {
			continue
		}
		return workGPU(data, target, devices)
	}
	return nil, errors.New("no gpu found")
}

func workGPU(data []byte, target uint64, devices []*cl.Device) (work []byte, err error) {
	work = make([]byte, 8)
	context, err := cl.CreateContext(devices)
	if err != nil {
		return
	}
	defer context.Release()
	program, err := context.CreateProgramWithSource([]string{kernelSource})
	if err != nil {
		return
	}
	defer program.Release()
	if err = program.BuildProgram(devices, ""); err != nil {
		return
	}
	kernel, err := program.CreateKernel("nano_work")
	if err != nil {
		return
	}
	defer kernel.Release()
	attempt, err := context.CreateEmptyBuffer(cl.MemReadOnly, len(work))
	if err != nil {
		return
	}
	defer attempt.Release()
	result, err := context.CreateBuffer(cl.MemWriteOnly|cl.MemCopyHostPtr, work)
	if err != nil {
		return
	}
	defer result.Release()
	root, err := context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, data)
	if err != nil {
		return
	}
	defer root.Release()
	if err = kernel.SetArgs(attempt, result, root, target); err != nil {
		return
	}
	queue, err := context.CreateCommandQueue(devices[0], 0)
	if err != nil {
		return
	}
	defer queue.Release()
	buf := make([]byte, len(work))
	for x := rand.Uint64(); bytes.Count(work, []byte{0}) == len(work); x += 1 << 20 {
		binary.LittleEndian.PutUint64(buf, x)
		if _, err = queue.EnqueueWriteBufferByte(attempt, false, 0, buf, nil); err != nil {
			return
		}
		if _, err = queue.EnqueueNDRangeKernel(kernel, nil, []int{1 << 20}, nil, nil); err != nil {
			return
		}
		if _, err = queue.EnqueueReadBufferByte(result, true, 0, work, nil); err != nil {
			return
		}
	}
	return
}

// GenerateCPU generates proof-of-work using the CPU.
func GenerateCPU(data []byte, target uint64) (work []byte, err error) {
	n := runtime.NumCPU()
	ch := make(chan []byte, n)
	hash := make([]hash.Hash, n)
	for i := 0; i < n; i++ {
		if hash[i], err = blake2b.New(8, nil); err != nil {
			return
		}
	}
	done := false
	x := rand.Uint64()
	for i := 0; i < n; i++ {
		go func(i int) {
			work := make([]byte, 8)
			for x := x + uint64(i); !done; x += uint64(n) {
				binary.BigEndian.PutUint64(work, x)
				hash[i].Reset()
				hash[i].Write(work)
				hash[i].Write(data)
				if binary.LittleEndian.Uint64(hash[i].Sum(nil)) >= target {
					done = true
					ch <- work
				}
			}
		}(i)
	}
	return <-ch, nil
}
