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

	var wrap = function (s) {
		var indent = "";
		var toknum = 0;
		for (var i = 0; i < s.length; i++) {
			indent += ' ';
			if (s[i] == ' ' && toknum == 0) {
				toknum++;
				continue;
			}
			if (s[i] != ' ' && toknum == 1) {
				break;
			}
		}

		var lines = [];
		var buf = "";
		var tokens = s.split(/[ \r\t\n]/);
		for (var i = 0; i < tokens.length; i++) {
			if (buf.length + tokens[i].length > 80) {
				lines.push(buf);
				buf = tokens[i];
			} else {
				buf += ' ' + tokens[i];
			}
		}
		if (buf.length > 0) {
			lines.push(buf);
		}
		return lines.join("<br/>" + indent + " <span class=\"cont\">&#8594;</span> ");
	}

	// colorize debugging output
	var colorize = function (s) {
		out = []
		lines = s.split("\n");
		for (i = 0; i < lines.length; i++) {
			var line = lines[i];
			var html = wrap(line);
			if (line.match(/^--+> data flow:/)) {
				out.push('<span class="data-flow">'+html+'</span>');
			} else if (line.match(/^--+>\s*$/)) {
				out.push('<span class="trace">&nbsp;</span>');
			} else if (line.match(/^--+>/)) {
				out.push('<span class="trace">'+html+'</span>');

			} else if (line.match(/^DEBUG> parsing/)) {
				out.push('<span class="debug debug-parsing">'+html+'</span>');
			} else if (line.match(/^DEBUG> running/)) {
				out.push('<span class="debug debug-running">'+html+'</span>');
			} else if (line == 'DEBUG> ') {
				out.push('<span class="debug">&nbsp;</span>');
			} else if (line.match(/^DEBUG>/)) {
				out.push('<span class="debug">'+html+'</span>');
			} else {
				out.push('<span>'+html+'</span>');
			}
		}
		return out.join("");
	}

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
			}),
			success: function(raw) {
				data = JSON.parse(raw);
				console.dir(data);

				$('#stdout, #stderr').hide();

				if (data.stdout != "") {
					$('#stdout').html(data.stdout).show();
				}
				if (data.stderr != "") {
					$('#stderr').html(colorize(data.stderr)).show();
				}
			}
		});
	})
})
