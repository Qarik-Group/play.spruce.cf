$(function () {
	// when `[ ] Trace` is checked,
	//      `[ ] Debugging` also needs to be on
	$('#trace').change(function (event) {
		if ($(event.target).is(':checked')) {
			$('#debug').prop('checked', true);
		}
	})

	// when `[ ] Debugging` is unchecked,
	//      `[ ] Trace` needs to be turned off
	$('#debug').change(function (event) {
		if (!$(event.target).is(':checked')) {
			$('#trace').prop('checked', false);
		}
	})

	// turn a prune list (text) into a list of parsed keys
	// (ignoring comments and blank lines).
	var prune = function (s) {
		var keys = [];
		lines = s.split("\n");
		for (i = 0; i < lines.length; i++) {
			line = lines[i].replace(/^\s+|\s+$|^\s*#.*$/, '');
			if (line != "") {
				keys.push(line)
			}
		}
		return keys;
	}

	// turn a multi-document YAML string into a list of
	// separate document strings.
	var yamls = function (s) {
		var docs = [];
		var doc = [];
		lines = s.split("\n");
		for (i = 0; i < lines.length; i++) {
			if (lines[i] == "---" && doc.length > 0) {
				docs.push(doc.join("\n") + "\n");
				doc = [];
			}
			doc.push(lines[i])
		}
		if (doc.length > 0) {
			docs.push(doc.join("\n") + "\n");
		}
		return docs;
	}

	// Set up the CodeMirror editor for syntax higlighting
	$('.yaml.code').each(function (i, textarea) {
		CodeMirror.fromTextArea(textarea, {
			theme       : 'spruce',
			lineNumbers : true
		});
	})

	// Handle [Merge] clicks and form submission
	$('#playground').submit(function (event) {
		event.preventDefault();
		$.ajax({
			type: 'POST',
			url:  '/spruce',
			data: JSON.stringify({
				yaml  : yamls($('#yaml').val()),
				debug : $('#debug').is(':checked'),
				trace : $('#trace').is(':checked'),
				prune : prune($('#prune').val())
			})
		});
	})
})
