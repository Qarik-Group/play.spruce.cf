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

	// turn a textarea into a list of parsed keys
	// (ignoring comments and blank lines).
	var text_to_list = function (s) {
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
		var n = 1;
		var docify = function (doc) {
			file = null;
			// find the first `# <something>.yml` comment
			for (var i = 0; i < doc.length; i++) {
				m = doc[i].match(/^---\s*#\s*(.+\.yml)\s*$/);
				if (m) {
					file = m[1];
					break;
				}
				m = doc[i].match(/^\s*#\s*(.+\.yml)\s*$/);
				if (m) {
					file = m[1];
					break;
				}
			}
			if (file == null) {
				file = "file" + n.toString() + ".yml"
			}
			n++;
			return { filename: file, contents: doc.join("\n") + "\n" };
		}

		var docs = [];
		var doc = [];
		lines = s.split("\n");
		for (var i = 0; i < lines.length; i++) {
			if (lines[i].match(/^---/) && doc.length > 0) {
				docs.push(docify(doc));
				doc = [];
			}
			doc.push(lines[i])
		}
		if (doc.length > 0) {
			docs.push(docify(doc));
		}
		return docs;
	}

	// Set up the CodeMirror editor for syntax higlighting
	$('.yaml.code').each(function (i, textarea) {
		$(textarea).data('editor',
			CodeMirror.fromTextArea(textarea, {
				theme       : 'spruce',
				lineNumbers : true
			}));
	});

	// JSONify the form data entered so far...
	formdata = function () {
		return {
			flavor          : $('[name="flavor"]:checked').val(),
			yaml            : yamls($('#yaml').data('editor').getValue()),
			debug           : $('#debug').is(':checked'),
			trace           : $('#trace').is(':checked'),
			fallback_append : $('#fallback-append').is(':checked'),
			skip_eval       : $('#skip-eval').is(':checked'),
			go_patch        : $('#go-patch').is(':checked'),
			prune           : text_to_list($('#prune').val()),
			env             : text_to_list($('#env').val()),
			cherry_pick     : text_to_list($('#cherry-pick').val()),

		};
	};

	// Form submission is a no-op
	$('#playground').submit(function (event) {
		event.preventDefault();
	});

	// Handle [Merge] clicks to perform AJAX merge request
	$('#playground #merge').click(function (event) {
		$.ajax({
			type: 'POST',
			url:  '/spruce',
			data: JSON.stringify(formdata()),
			success: function(raw) {
				data = JSON.parse(raw);
				$('#about, #stdout, #stderr').hide();

				if (data.about  != "") { $('#about').html(data.about).show(); }
				if (data.stderr != "") { $('#stderr').html(colorize(data.stderr)).show(); }
				if (data.stdout != "") {
					$('#stdout').html('<textarea class="yaml code">'+data.stdout+'</textarea>')
					            .show();
					$('#stdout .yaml.code').each(function (i, textarea) {
						$(textarea).data('editor',
							CodeMirror.fromTextArea(textarea, {
								readOnly    : true, /* crucial! */
								theme       : 'spruce',
								lineNumbers : true
							}));
					});
				}
			}
		});
	});

	// Handle [Share] clicks to handle storage
	$('#playground #share').click(function (event) {
		event.preventDefault();
		$.ajax({
			type: 'POST',
			url:  '/mem',
			data: JSON.stringify(formdata()),
			success: function(key) {
				document.location.hash = key;
			}
		});
	});

	$.ajax({
		type: 'GET',
		url:  '/meta',
		success: function(raw) {
			data = JSON.parse(raw);

			$('#flavors').append('<li><input type="radio" value="' + data.flavors[0] + '" name="flavor" /> latest (' + data.flavors[0] + ')</li>');
			for (var i = 0; i < data.flavors.length; i++) {
				$('#flavors').append('<li><input type="radio" value="' + data.flavors[i] + '" name="flavor" /> ' + data.flavors[i] + '</li>');
			}
			$('#flavors li:first-child input').prop('checked', true);
		}
	});

	if (document.location.hash != "") {
		key = document.location.hash.replace(/^#/, '')
		$.ajax({
			type: 'GET',
			url:  '/mem?k=' + key,
			success: function(raw) {
				data = JSON.parse(raw);
				data.prune       = data.prune == null ? [] : data.prune;
				data.env         = data.env == null ? [] : data.env;
				data.cherry_pick = data.cherry_pick == null ? [] : data.cherry_pick;

				// keys to prune
				$('#prune').val("# one per line\n" + data.prune.join("\n") + "\n");
				// environment variables
				$('#env').val("# VAR=value\n" + data.env.join("\n") + "\n");
				// cherry-picks
				$('#cherry-pick').val("# one per line\n" + data.cherry_pick.join("\n") + "\n");

				// YAML document(s)
				var yaml = "";
				for (var i = 0; i < data.yaml.length; i++) {
					yaml += data.yaml[i].contents + "\n\n";
				}
				$('#yaml').data('editor').setValue(yaml);

				// Debug / Trace / go-patch/ skip-eval/ fallback-append flags
				if (data.trace)           { $('#trace').prop('checked', true); }
				if (data.debug)           { $('#debug').prop('checked', true); }
				if (data.go_patch)        { $('#go-patch').prop('checked', true); }
				if (data.skip_eval)       { $('#skip-eval').prop('checked', true); }
				if (data.fallback_append) { $('#fallback-append').prop('checked', true); }

				// Flavor
				$('[name=flavor][value="' + data.flavor + '"]').prop('checked', true);

				// Do a merge
				$('#merge').trigger('click');
			}
		});
	}
})
