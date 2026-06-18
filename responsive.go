package cloudinary_gem

import "sort"

type SrcSetBreakpoint struct {
	Width          int
	Transformation URLOptions
}

func SrcSet(cld *Cloudinary, publicID string, breakpoints []SrcSetBreakpoint, opts URLOptions) (string, error) {
	sort.Slice(breakpoints, func(i, j int) bool { return breakpoints[i].Width < breakpoints[j].Width })
	parts := []string{}
	for _, bp := range breakpoints {
		o := opts
		if bp.Transformation.Width != 0 {
			o.Width = bp.Transformation.Width
		} else {
			o.Width = bp.Width
		}
		if bp.Transformation.Height != 0 {
			o.Height = bp.Transformation.Height
		}
		if bp.Transformation.Crop != "" {
			o.Crop = bp.Transformation.Crop
		}
		u, err := cld.URL(publicID, o)
		if err != nil {
			return "", err
		}
		parts = append(parts, u+" "+intString(bp.Width)+"w")
	}
	return BuildArray(parts), nil
}

func SizesAttribute(widths []int) string {
	parts := []string{}
	for _, w := range widths {
		parts = append(parts, "(max-width: "+intString(w)+"px) "+intString(w)+"px")
	}
	return BuildArray(parts)
}
