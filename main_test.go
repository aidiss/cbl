package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateMineralType(t *testing.T) {
	assert := assert.New(t)
	mt := CreateMineralType("topaz", 4, 8, 8)
	assert.Equal("topaz", mt.Name, "ff")
	assert.Equal("topaz", mt.Name, "ff")
}

func TestCreateMineral(t *testing.T) {
	assert := assert.New(t)
	mt := CreateMineralType("topaz", 4, 8, 8)
	m := CreateMineral(mt, "fractured", 4)

	assert.Equal("fractured", m.State, "ff")

}

func TestNewJob(t *testing.T) {
	assert := assert.New(t)
	mt := CreateMineralType("topaz", 4, 8, 8)
	m := CreateMineral(mt, "fractured", 4)

	j := NewJob("fracture", &m, "NEW")
	assert.Equal("NEW", j.Status, "ff")
}

func TestFactory_FractureMineral(t *testing.T) {
	assert := assert.New(t)
	f := Factory{
		jobQueue:   nil,
		currentJob: Job{},
		active:     false,
	}
	mt := CreateMineralType("topaz", 4, 8, 8)
	m := CreateMineral(mt, "fractured", 8)

	err := f.FractureMineral(&m)
	assert.NotNil(err, "Problem")
}