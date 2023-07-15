package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	// Create a slice of anonymous structs containing the test case name, 
	// input to our humanDate() function (the tm field), and expected 
	// output (the want field).

	tests := []struct { 
		name string 
		tm time.Time 
		want string 
	}{ 
		{ name: "UTC", tm: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC), want: "17 Mar 2022 at 10:15", }, 
		{ name: "Empty", tm: time.Time{}, want: "", }, 
		{ name: "CET", tm: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)), want: "17 Mar 2022 at 09:15", },
	}

	// Loop over test cases
	for _, tt := range tests {
			// Run all time.Date objects through humanDate
			hd := humanDate(tt.tm)
		
			// Check that the output from the humanDate function is in the format we 
			// expect. If it isn't what we expect, use the t.Errorf() function to 
			// indicate that the test has failed and log the expected and actual values.
			if hd != tt.want { 
				t.Errorf("got %q; want %q", hd, tt.want) 
			}
	}
}