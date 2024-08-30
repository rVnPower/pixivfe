# About Testing

Do we really need testing? What to test?

## The Current State Of Testing

- some playwright test
- syntax checking jet templates

## Property testing with [flyingmutant/rapid](https://github.com/flyingmutant/rapid/)

I tried doing using rapid to test templates. However, the generator is broken. The error says `reflect.Set on unexported field`.

See the `rapid` branch.

`interface{}` can't be used anywhere in Data_* or else rapid will complain.
