package template

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = DescribeTable("Generate Upstream UID",
	func(prefix, input, output string) {
		g := &upstreamTemplateGenerator{
			upstream: &v1.Upstream{
				Metadata: &core.Metadata{
					Name: input,
				},
			},
			dashboardPrefix: prefix,
		}
		Expect(g.GenerateUid()).To(Equal(output))
	},
	Entry("Input does not require modification", "", "foo", "foo"),
	Entry("Input is greater than 40 characters", "", "prefix-0123456789012345678901234567890123456789", "0123456789012345678901234567890123456789"),
	Entry("Input has illegal characters", "", "foo.bar-biz!baz", "foo_bar_biz_baz"),
	Entry("Input has a prefix", "prefix-", "foo", "prefix_foo"),
)
