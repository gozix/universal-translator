// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configurator

import (
	ut "github.com/go-playground/universal-translator"
)

// Configurator configures universal translator after its creation.
type Configurator func(*ut.UniversalTranslator) error
