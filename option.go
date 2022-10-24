// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package ut

import "github.com/go-playground/locales"

type (
	// Option interface.
	Option interface {
		apply(b *Bundle)
	}

	// optionFunc wraps a func so it satisfies the Option interface.
	optionFunc func(b *Bundle)
)

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
