package check_in

import "testing"

func TestReturnResultSuccess(t *testing.T) {
	got := returnResult(map[string]interface{}{
		"error_code": float64(0),
		"data": map[string]interface{}{
			"continue_checkin_days": float64(7),
		},
	})
	want := "恭喜签到成功！您已连续签到7天!"
	if got != want {
		t.Fatalf("returnResult()=%q, want %q", got, want)
	}
}
