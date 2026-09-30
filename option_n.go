package safe

// Lots of boilerplate glue to map Option<n> alias types onto the concrete optionN implementation
// We only go up to Option8; if you need more then you should probably reconsider your life choices

type Option2[T1, T2 any] optionN[T1, T2, xx, xx, xx, xx, xx, xx]
type Option3[T1, T2, T3 any] optionN[T1, T2, T3, xx, xx, xx, xx, xx]
type Option4[T1, T2, T3, T4 any] optionN[T1, T2, T3, T4, xx, xx, xx, xx]
type Option5[T1, T2, T3, T4, T5 any] optionN[T1, T2, T3, T4, T5, xx, xx, xx]
type Option6[T1, T2, T3, T4, T5, T6 any] optionN[T1, T2, T3, T4, T5, T6, xx, xx]
type Option7[T1, T2, T3, T4, T5, T6, T7 any] optionN[T1, T2, T3, T4, T5, T6, T7, xx]
type Option8[T1, T2, T3, T4, T5, T6, T7, T8 any] optionN[T1, T2, T3, T4, T5, T6, T7, T8]

func None2[T1, T2 any]() Option2[T1, T2] {
	return Option2[T1, T2]{}
}
func None3[T1, T2, T3 any]() Option3[T1, T2, T3] {
	return Option3[T1, T2, T3]{}
}
func None4[T1, T2, T3, T4 any]() Option4[T1, T2, T3, T4] {
	return Option4[T1, T2, T3, T4]{}
}
func None5[T1, T2, T3, T4, T5 any]() Option5[T1, T2, T3, T4, T5] {
	return Option5[T1, T2, T3, T4, T5]{}
}
func None6[T1, T2, T3, T4, T5, T6 any]() Option6[T1, T2, T3, T4, T5, T6] {
	return Option6[T1, T2, T3, T4, T5, T6]{}
}
func None7[T1, T2, T3, T4, T5, T6, T7 any]() Option7[T1, T2, T3, T4, T5, T6, T7] {
	return Option7[T1, T2, T3, T4, T5, T6, T7]{}
}
func None8[T1, T2, T3, T4, T5, T6, T7, T8 any]() Option8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Option8[T1, T2, T3, T4, T5, T6, T7, T8]{}
}

func Some2[T0, T1, T2 any](value T0) Option2[T1, T2] {
	return (Option2[T1, T2])(someN[T0, T1, T2, xx, xx, xx, xx, xx, xx](value))
}
func Some3[T0, T1, T2, T3 any](value T0) Option3[T1, T2, T3] {
	return (Option3[T1, T2, T3])(someN[T0, T1, T2, T3, xx, xx, xx, xx, xx](value))
}
func Some4[T0, T1, T2, T3, T4 any](value T0) Option4[T1, T2, T3, T4] {
	return (Option4[T1, T2, T3, T4])(someN[T0, T1, T2, T3, T4, xx, xx, xx, xx](value))
}
func Some5[T0, T1, T2, T3, T4, T5 any](value T0) Option5[T1, T2, T3, T4, T5] {
	return (Option5[T1, T2, T3, T4, T5])(someN[T0, T1, T2, T3, T4, T5, xx, xx, xx](value))
}
func Some6[T0, T1, T2, T3, T4, T5, T6 any](value T0) Option6[T1, T2, T3, T4, T5, T6] {
	return (Option6[T1, T2, T3, T4, T5, T6])(someN[T0, T1, T2, T3, T4, T5, T6, xx, xx](value))
}
func Some7[T0, T1, T2, T3, T4, T5, T6, T7 any](value T0) Option7[T1, T2, T3, T4, T5, T6, T7] {
	return (Option7[T1, T2, T3, T4, T5, T6, T7])(someN[T0, T1, T2, T3, T4, T5, T6, T7, xx](value))
}
func Some8[T0, T1, T2, T3, T4, T5, T6, T7, T8 any](value T0) Option8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return (Option8[T1, T2, T3, T4, T5, T6, T7, T8])(someN[T0, T1, T2, T3, T4, T5, T6, T7, T8](value))
}

func (o Option2[T1, T2]) Match() any {
	return (optionN[T1, T2, xx, xx, xx, xx, xx, xx])(o).match()
}
func (o Option3[T1, T2, T3]) Match() any {
	return (optionN[T1, T2, T3, xx, xx, xx, xx, xx])(o).match()
}
func (o Option4[T1, T2, T3, T4]) Match() any {
	return (optionN[T1, T2, T3, T4, xx, xx, xx, xx])(o).match()
}
func (o Option5[T1, T2, T3, T4, T5]) Match() any {
	return (optionN[T1, T2, T3, T4, T5, xx, xx, xx])(o).match()
}
func (o Option6[T1, T2, T3, T4, T5, T6]) Match() any {
	return (optionN[T1, T2, T3, T4, T5, T6, xx, xx])(o).match()
}
func (o Option7[T1, T2, T3, T4, T5, T6, T7]) Match() any {
	return (optionN[T1, T2, T3, T4, T5, T6, T7, xx])(o).match()
}
func (o Option8[T1, T2, T3, T4, T5, T6, T7, T8]) Match() any {
	return (optionN[T1, T2, T3, T4, T5, T6, T7, T8])(o).match()
}

func (o Option2[T1, T2]) IsSome() bool {
	return o.t != nil
}
func (o Option3[T1, T2, T3]) IsSome() bool {
	return o.t != nil
}
func (o Option4[T1, T2, T3, T4]) IsSome() bool {
	return o.t != nil
}
func (o Option5[T1, T2, T3, T4, T5]) IsSome() bool {
	return o.t != nil
}
func (o Option6[T1, T2, T3, T4, T5, T6]) IsSome() bool {
	return o.t != nil
}
func (o Option7[T1, T2, T3, T4, T5, T6, T7]) IsSome() bool {
	return o.t != nil
}
func (o Option8[T1, T2, T3, T4, T5, T6, T7, T8]) IsSome() bool {
	return o.t != nil
}

func (o Option2[T1, T2]) IsNone() bool {
	return o.t == nil
}
func (o Option3[T1, T2, T3]) IsNone() bool {
	return o.t == nil
}
func (o Option4[T1, T2, T3, T4]) IsNone() bool {
	return o.t == nil
}
func (o Option5[T1, T2, T3, T4, T5]) IsNone() bool {
	return o.t == nil
}
func (o Option6[T1, T2, T3, T4, T5, T6]) IsNone() bool {
	return o.t == nil
}
func (o Option7[T1, T2, T3, T4, T5, T6, T7]) IsNone() bool {
	return o.t == nil
}
func (o Option8[T1, T2, T3, T4, T5, T6, T7, T8]) IsNone() bool {
	return o.t == nil
}
