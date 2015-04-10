package main

import (
	"reflect"
	"testing"
)

func TestParseExpr(t *testing.T) {
	for _, tc := range []struct {
		in    string
		val   Value
		isErr bool
	}{
		{
			in:  "foo",
			val: literal("foo"),
		},
		{
			in:  "(foo)",
			val: literal("(foo)"),
		},
		{
			in:  "{foo}",
			val: literal("{foo}"),
		},
		{
			in:  "$$",
			val: literal("$"),
		},
		{
			in:  "foo$$bar",
			val: literal("foo$bar"),
		},
		{
			in:  "$foo",
			val: Expr{varref{varname: literal("f")}, literal("oo")},
		},
		{
			in:  "$(foo)",
			val: varref{varname: literal("foo")},
		},
		{
			in: "$(foo:.c=.o)",
			val: varsubst{
				varname: literal("foo"),
				pat:     literal(".c"),
				subst:   literal(".o"),
			},
		},
		{
			in: "$(subst $(space),$(,),$(foo))/bar",
			val: Expr{
				&funcSubst{
					fclosure: fclosure{
						args: []Value{
							varref{
								varname: literal("space"),
							},
							varref{
								varname: literal(","),
							},
							varref{
								varname: literal("foo"),
							},
						},
					},
				},
				literal("/bar"),
			},
		},
		{
			in: "$(subst $(space),$,,$(foo))",
			val: &funcSubst{
				fclosure: fclosure{
					args: []Value{
						varref{
							varname: literal("space"),
						},
						varref{
							varname: literal(""),
						},
						Expr{
							literal(","),
							varref{
								varname: literal("foo"),
							},
						},
					},
				},
			},
		},
		{
			in: `$(shell echo '()')`,
			val: &funcShell{
				fclosure: fclosure{
					args: []Value{
						literal("echo '()'"),
					},
				},
			},
		},
		{
			in: `$(shell echo '()')`,
			val: &funcShell{
				fclosure: fclosure{
					args: []Value{
						literal("echo '()'"),
					},
				},
			},
		},
		{
			in: `$(shell echo ')')`,
			val: Expr{
				&funcShell{
					fclosure: fclosure{
						args: []Value{
							literal("echo '"),
						},
					},
				},
				literal("')"),
			},
		},
		{
			in: `${shell echo ')'}`,
			val: &funcShell{
				fclosure: fclosure{
					args: []Value{
						literal("echo ')'"),
					},
				},
			},
		},
		{
			in: `${shell echo '}'}`,
			val: Expr{
				&funcShell{
					fclosure: fclosure{
						args: []Value{
							literal("echo '"),
						},
					},
				},
				literal("'}"),
			},
		},
		{
			in: `$(shell make --version | ruby -n0e 'puts $$_[/Make (\d)/,1]')`,
			val: &funcShell{
				fclosure: fclosure{
					args: []Value{
						literal(`make --version | ruby -n0e 'puts $_[/Make (\d)/,1]'`),
					},
				},
			},
		},
	} {
		val, _, err := parseExpr([]byte(tc.in), nil)
		if tc.isErr {
			if err == nil {
				t.Errorf(`parseExpr(%q)=_, _, nil; want error`, tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf(`parseExpr(%q)=_, _, %v; want nil error`, tc.in, err)
			continue
		}
		if got, want := val, tc.val; !reflect.DeepEqual(got, want) {
			t.Errorf("parseExpr(%[1]q)=%[2]q %#[2]v, _, _;\n want %[3]q %#[3]v, _, _", tc.in, got, want)
		}
	}
}
