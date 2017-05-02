package handlers

import "syscall"

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

func GetDiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func DiskInfoByMeasure(path string, measure uint64) DiskStatus {
	d := GetDiskUsage(path)
	d.All = uint64(float64(d.All) / float64(measure))
	d.Used = uint64(float64(d.Used) / float64(measure))
	d.Free = uint64(float64(d.Free) / float64(measure))
	return d
}
