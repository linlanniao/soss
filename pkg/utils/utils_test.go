package utils_test

import (
	"reflect"
	"testing"

	"github.com/linlanniao/soss/pkg/utils"
)

func TestRemoveDuplicates(t *testing.T) {

	intCases := []struct {
		name string
		args []int
		want []int
	}{
		{
			name: "remove duplicates from int slice",
			args: []int{1, 1, 2, 2, 3, 3},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range intCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.RemoveDuplicates(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
	float32Cases := []struct {
		name string
		args []float32
		want []float32
	}{
		{
			name: "remove duplicates from float32 slice",
			args: []float32{1, 1, 2, 2, 3, 3},
			want: []float32{1, 2, 3},
		},
	}

	for _, tt := range float32Cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.RemoveDuplicates(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}

	stringCases := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "remove duplicates from string slice",
			args: []string{"a", "a", "b", "b", "c", "c"},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range stringCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.RemoveDuplicates(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}

	type s struct {
		name string
		age  int64
		sex  string
	}

	structCases := []struct {
		name string
		args []s
		want []s
	}{
		{
			name: "remove duplicates from struct slice",
			args: []s{
				{
					name: "a",
					age:  20,
					sex:  "male",
				},
				{
					name: "a",
					age:  20,
					sex:  "male",
				},
				{
					name: "b",
					age:  25,
					sex:  "female",
				},
				{
					name: "b",
					age:  25,
					sex:  "female",
				},
				{
					name: "c",
					age:  22,
					sex:  "male",
				},
			},
			want: []s{
				{
					name: "a",
					age:  20,
					sex:  "male",
				},
				{
					name: "b",
					age:  25,
					sex:  "female",
				},
				{
					name: "c",
					age:  22,
					sex:  "male",
				},
			},
		},
	}
	for _, tt := range structCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.RemoveDuplicates(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
