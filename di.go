// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package ut

import "github.com/gozix/di"

const (
	// TagTranslator is a tag for marking locale translator without overriding existing.
	TagTranslator = "universal-translator.locale-translator"

	// TagTranslatorOverride is a tag for marking locale translator with overriding existing.
	TagTranslatorOverride = "universal-translator.locale-translator.override"
)

func AsTranslator(override bool) di.Tags {
	if override {
		return di.Tags{{
			Name: TagTranslatorOverride,
		}}
	}

	return di.Tags{{
		Name: TagTranslator,
	}}
}

func withTranslator(override bool) di.Modifier {
	if override {
		return di.WithTags(TagTranslatorOverride)
	}

	return di.WithTags(TagTranslator)
}
