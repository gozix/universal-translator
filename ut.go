package ut

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/sarulabs/di/v2"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct {
		fallback locales.Translator
		locales  []locales.Translator
	}

	// Configurator is configurator func interface.
	Configurator = func(*ut.UniversalTranslator) error

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

	// TranslatorConfiguratorName is definition name.
	TranslatorConfiguratorName = "universal-translator.configurator.translator"

	// TagConfigurator is tag to mark configurator injections.
	TagConfigurator = "universal-translator.configurator"

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
	return builder.Add(
		di.Def{
			Name: BundleName,
			Build: func(ctn di.Container) (_ interface{}, err error) {
				var configurators = make([]Configurator, 0, 4)
				for name, def := range ctn.Definitions() {
					for _, tag := range def.Tags {
						if TagConfigurator != tag.Name {
							continue
						}

						var configurator Configurator
						if err = ctn.Fill(name, &configurator); err != nil {
							return nil, err
						}

						configurators = append(configurators, configurator)
					}
				}

				var translator = ut.New(b.fallback, b.locales...)
				for _, configurator := range configurators {
					if err = configurator(translator); err != nil {
						return nil, err
					}
				}

				return translator, nil
			},
		}, di.Def{
			Name: TranslatorConfiguratorName,
			Build: func(ctn di.Container) (interface{}, error) {
				return func(translator ut.UniversalTranslator) (err error) {
					for name, def := range ctn.Definitions() {
						for _, tag := range def.Tags {
							if tag.Name != TagTranslator {
								continue
							}

							var localeTranslator locales.Translator
							if err = ctn.Fill(name, &localeTranslator); err != nil {
								return err
							}

							_, override := tag.Args[TagArgOverride]
							if err = translator.AddTranslator(localeTranslator, override); err != nil {
								return err
							}

							break
						}
					}

					return nil
				}, nil
			},
		},
	)
}

// apply implements Option.
func (f optionFunc) apply(bundle *Bundle) {
	f(bundle)
}
