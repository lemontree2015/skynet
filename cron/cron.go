package cron

import (
	"time"
)

// 备注:
// 之类所有执行的函数, 因为都是在独立的goroutine中执行, 所以都需要recover保护.

// 在一个特定的时间点, 执行fn
type CronAt struct {
	timer   *time.Timer
	fn      func()
	is_stop bool
}

func NewCronAt(t time.Time, fn func()) *CronAt {
	cron := &CronAt{
		timer:   nil,
		fn:      fn,
		is_stop: false,
	}

	registerAt(cron, t, fn)

	return cron
}

func registerAt(cron *CronAt, t time.Time, fn func()) {
	cron.timer = time.AfterFunc(t.Sub(time.Now()), func() {
		if !cron.is_stop {
			fn()
		}
		cron.is_stop = true
	})
}

func (cron *CronAt) IsStopped() bool {
	return cron.is_stop
}

func (cron *CronAt) Stop() bool {
	cron.is_stop = true
	if cron.timer != nil {
		return cron.timer.Stop()
	}

	return false
}

// 经过duration之后, 再执行fn
type CronAfter struct {
	timer   *time.Timer
	fn      func()
	is_stop bool
}

func NewCronAfter(duration time.Duration, fn func()) *CronAfter {
	cron := &CronAfter{
		timer:   nil,
		fn:      fn,
		is_stop: false,
	}

	registerAfter(cron, duration, fn)

	return cron
}

func registerAfter(cron *CronAfter, duration time.Duration, fn func()) {
	cron.timer = time.AfterFunc(duration, func() {
		if !cron.is_stop {
			fn()
		}
		cron.is_stop = true
	})
}

func (cron *CronAfter) IsStopped() bool {
	return cron.is_stop
}

func (cron *CronAfter) Stop() bool {
	cron.is_stop = true
	if cron.timer != nil {
		return cron.timer.Stop()
	}

	return false
}

func (cron *CronAfter) Reset(duration time.Duration) bool {
	if cron.timer != nil {
		return cron.timer.Reset(duration)
	}
	return false
}

// 每间隔duration, 执行一次fn
type CronEvery struct {
	timer   *time.Timer
	fn      func()
	is_stop bool
}

func NewCronEvery(duration time.Duration, fn func()) *CronEvery {
	cron := &CronEvery{
		timer:   nil,
		fn:      fn,
		is_stop: false,
	}

	registerEvery(cron, duration, fn)

	return cron
}

// 每间隔duration, 执行一次fn
func registerEvery(cron *CronEvery, duration time.Duration, fn func()) {
	cron.timer = time.AfterFunc(duration, func() {
		if !cron.is_stop {
			fn()
		}

		if !cron.is_stop {
			registerEvery(cron, duration, fn)
		}
	})
}

func (cron *CronEvery) IsStopped() bool {
	return cron.is_stop
}

func (cron *CronEvery) Stop() bool {
	cron.is_stop = true
	if cron.timer != nil {
		return cron.timer.Stop()
	}

	return false
}

type CronUntil struct {
	timer   *time.Timer
	fn      func()
	is_stop bool
}

// 每间隔duration, 执行一次fn, 直到时间到达t结束
func NewCronUntil(t time.Time, duration time.Duration, fn func()) *CronUntil {
	cron := &CronUntil{
		timer:   nil,
		fn:      fn,
		is_stop: false,
	}

	registerUntil(cron, t, duration, fn)

	return cron
}
func registerUntil(cron *CronUntil, t time.Time, duration time.Duration, fn func()) {
	if t.Sub(time.Now()) > 0 {
		cron.timer = time.AfterFunc(duration, func() {
			if !cron.is_stop {
				fn()
			}

			if !cron.is_stop {
				registerUntil(cron, t, duration, fn)
			}
		})
		return
	}

	cron.is_stop = true
}

func (cron *CronUntil) IsStopped() bool {
	return cron.is_stop
}

func (cron *CronUntil) Stop() bool {
	cron.is_stop = true
	if cron.timer != nil {
		return cron.timer.Stop()
	}

	return false
}
