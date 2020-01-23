/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openebs

import (
	"github.com/mayadata-io/openebs-operator/types"
)

// Set the default values for helpers used by OpenEBS components
// such as linux-utils, etc.
func (r *Reconciler) setHelperDefaultsIfNotSet() error {
	if r.OpenEBS.Spec.Helper == nil {
		r.OpenEBS.Spec.Helper = &types.Helper{}
	}
	// form the linux-utils image
	if r.OpenEBS.Spec.Helper.ImageTag == "" {
		r.OpenEBS.Spec.Helper.ImageTag = r.OpenEBS.Spec.Version
	}
	r.OpenEBS.Spec.Helper.Image = r.OpenEBS.Spec.ImagePrefix +
		"linux-utils:" + r.OpenEBS.Spec.Helper.ImageTag
	return nil
}
