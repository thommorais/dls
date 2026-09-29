/**
 * A run of evenly spaced grays, darkest first, so stacked series stay
 * tellable apart without reaching for hue. Used wherever a chart used to
 * lean on the API's per-type color.
 */
export function grayShades(count: number): readonly string[] {
	if (count <= 1) return ['hsl(0 0% 20%)']
	return Array.from({ length: count }, (_, index) => {
		const t = index / (count - 1)
		return `hsl(0 0% ${18 + t * 55}%)`
	})
}
