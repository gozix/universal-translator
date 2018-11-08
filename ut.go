package ut

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/sarulabs/di"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct {
		fallback locales.Translator
		locales  []locales.Translator
	}

	// Translator type alias of ut.Translator.
	Translator = ut.Translator

	// UniversalTranslator type alias of ut.UniversalTranslator.
	UniversalTranslator = ut.UniversalTranslator

	// Option interface.
	Option interface {
		apply(b *Bundle)
	}

	// optionFunc wraps a func so it satisfies the Option interface.
	optionFunc func(b *Bundle)
)

const (
	// BundleName is default definition name.
	BundleName = "universal-translator"

	// TagTranslator is tag to mark injected locale translator.
	TagTranslator = "universal-translator.locale-translator"

	// TagArgOverride is tag argument name to override locale translator.
	TagArgOverride = "override"
)

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

// Fallback option.
func Fallback(fallback locales.Translator) Option {
	return optionFunc(func(b *Bundle) {
		b.fallback = fallback
		b.locales = append(b.locales, fallback)
	})
}

// Locales option.
func Locales(locales ...locales.Translator) Option {
	return optionFunc(func(b *Bundle) {
		b.locales = append(locales, locales...)
	})
}

// Key implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	builder.Add(di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var translator = ut.New(b.fallback, b.locales...)
			for name, def := range ctn.Definitions() {
				for _, tag := range def.Tags {
					if tag.Name != TagTranslator {
						continue
					}

					var localeTranslator locales.Translator
					if err = ctn.Fill(name, &localeTranslator); err != nil {
						return nil, err
					}

					_, override := tag.Args[TagArgOverride]
					if err = translator.AddTranslator(localeTranslator, override); err != nil {
						return nil, err
					}

					break
				}
			}

			return translator, nil
		},
	})

	return nil
}

// apply implements Option.
func (f optionFunc) apply(bundle *Bundle) {
	f(bundle)
}
