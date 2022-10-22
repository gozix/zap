// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package zap

import "github.com/gozix/di"

// tagCoreFactory is factory tag.
const tagCoreFactory = "zap.core.factory"

// AsCoreFactory is syntax sugar for the di container.
func AsCoreFactory() di.ProvideOption {
	return di.ProvideOptions(
		di.As(new(CoreFactory)),
		di.Tags{{
			Name: tagCoreFactory,
		}},
	)
}

func withCoreFactory() di.Modifier {
	return di.WithTags(tagCoreFactory)
}
