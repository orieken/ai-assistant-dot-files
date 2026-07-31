# Eval Rubric: accessibility-engineer / input-component.vue

- **Missing label is identified with a concrete fix**: the `<input>` has no `<label>` or `aria-label`. Agent identifies this as a WCAG 1.3.1 (Level A) violation and provides the exact fix (either a visible `<label for="...">` or `aria-label` on the input) — not just "add a label."
- **Icon div flagged as non-interactive element**: the `<div class="icon">` with a click handler is not keyboard-accessible and has no `role`. Agent recommends replacing with `<button>` (semantically interactive, keyboard-focusable by default) — not just adding `role="button"` without `tabindex="0"`.
- **Placeholder color contrast flagged**: `color: #999` on `background: #eee` fails WCAG 1.4.3 (contrast ratio ~2.8:1, minimum is 4.5:1 for normal text). Agent names the failing ratio and suggests a specific corrected hex value.
- **Error message needs aria-live**: the `errorMsg` div has no `aria-live="polite"` — screen readers won't announce dynamically injected error text. Agent identifies this and the correct `aria-live` value.
- **List items missing keyboard interaction**: `<li>` elements with `@click` handlers are not keyboard-accessible. Agent recommends converting to `<button>` elements or adding `role="option"` within an `aria-listbox` pattern with keyboard event handling.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
