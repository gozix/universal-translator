// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package ut

import (
	"github.com/gozix/di"
	"github.com/gozix/glue/v3"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

// Bundle implements the glue.Bundle interface.
type Bundle struct {
	fallback locales.Translator
	locales  []locales.Translator
}

// Bundle implements the glue.Bundle interface.
var _ glue.Bundle = (*Bundle)(nil)

// BundleName is default definition name.
const BundleName = "universal-translator"

// NewBundle create bundle instance.
func NewBundle(options ...Option) *Bundle {
	var (
		locale = en.New()
		bundle = Bundle{
			fallback: locale,
			locales:  []locales.Translator{locale},
		}
	)

	for _, option := range options {
		option.apply(&bundle)
	}

	return &bundle
}

func (b *Bundle) Name() string {
	return BundleName
}

func (b *Bundle) Build(builder di.Builder) error {
	return builder.Provide(
		b.provideUT,
		di.Constraint(0, di.Optional(true), withTranslator(false)),
		di.Constraint(1, di.Optional(true), withTranslator(true)),
	)
}

func (b *Bundle) provideUT(append []locales.Translator, override []locales.Translator) (_ *ut.UniversalTranslator, err error) {
	var translator = ut.New(b.fallback, b.locales...)

	for _, localeTranslator := range append {
		if err = translator.AddTranslator(localeTranslator, false); err != nil {
			return nil, err
		}
	}

	for _, localeTranslator := range override {
		if err = translator.AddTranslator(localeTranslator, true); err != nil {
			return nil, err
		}
	}

	return translator, nil
}

// apply implements Option.
func (f optionFunc) apply(bundle *Bundle) {
	f(bundle)
}
