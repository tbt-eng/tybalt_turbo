<!-- The red asterisk after a required field's label. It is hidden from screen
readers: the control says it is required with aria-required instead. The DS
inputs render it when given `required`; a label written by hand puts it inside
the <label>, after the text. DSRequiredLegend explains it. The asterisk is CSS
content, so the label's text stays the field's name for tests and searches. -->
<span class="ml-0.5 text-red-600 after:content-['*']" aria-hidden="true"></span>
