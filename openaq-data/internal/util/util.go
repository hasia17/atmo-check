package util

import "time"

func WaitFor(ch ...chan struct{}) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for _, c := range ch {
			<-c
		}
		close(done)
	}()
	return done
}

func RemoveDuplicates[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	var result []T
	for _, item := range input {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func StringToTime(dateStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		t, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return time.Time{}, err
		}
	}
	return t, nil
}
