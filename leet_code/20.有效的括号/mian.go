package main

func getTopValue(stack []byte) byte{
	if len(stack) == 0 {
		return 0
	}

	return stack[len(stack) - 1]
}

func isValid(s string) bool {
	stack := make([]byte, 0)

	for _, v := range s {
		switch v {
		case '{', '[', '(':
			stack = append(stack, byte(v))
		case '}':
			if topValue := getTopValue(stack); topValue == '{' {
				stack = stack[:len(stack) - 1]
			} else {
				return false
			}
		case ']':
			if topValue := getTopValue(stack); topValue == '[' {
				stack = stack[:len(stack) - 1]
			} else {
				return false
			}
		case ')':
			if topValue := getTopValue(stack); topValue == '(' {
				stack = stack[:len(stack) - 1]
			} else {
				return false
			}
		default:
			return false
		}
	}

	if len(stack) == 0 {
		return true
	}

	return false
}

func main() {

}
