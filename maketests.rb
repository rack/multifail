# TODO(raggi): reach out to masinter@adobe.com and plh@w3.org and get this party started.

# TODO(raggi): add tests that exercise the FormData API.
def render(testname, page_charset, form_charset, fields)
  open("tests/#{testname}.html", "w") do |io|
    if page_charset && encoding = Encoding.find(page_charset)
      io.set_encoding(encoding) if encoding
      begin
        fields.each { |f| f.encode(encoding) }
      rescue Encoding::UndefinedConversionError
        $stderr.puts "Cannot encode fields for #{testname} in #{encoding}"
        next
      end
    end
    etestname = page_charset ? testname.encode(page_charset, :xml => :text) : testname
    quote = testname.include?('"') ? "'" : '"'
    io.write <<-HTML
<!doctype html>
<html>
<head>
<title>#{etestname}</title>
#{"<meta charset='#{page_charset}'>" if page_charset}
</head>
<body>
<form action=#{quote}#{etestname}#{quote} enctype="multipart/form-data" #{"accept-charset='#{form_charset}'" if form_charset}>
  #{fields.join("\n  ")}
  <input type="submit">
</form>
</body>
</html>
    HTML
  end
end


charsets = %w[utf-8 iso-8859-1 gb2312 windows-1251 windows-1252 shift_jis gbk] + [nil]

fieldsets = {
  empty_text: [{name: 'text', type: 'text', value: nil}],
  with_text:  [{name: 'text', type: 'text', value: 'text'}],
  empty_file: [{name: 'file', type: 'file', value: nil}],
  with_file:  [{name: 'file', type: 'file', value: File}],

  two_empty_text: [
    {name: 'text', type: 'text', value: nil},
    {name: 'text2', type: 'text', value: nil},
  ],
  two_with_text: [
    {name: 'text', type: 'text', value: 'text'},
    {name: 'text2', type: 'text', value: 'text2'}
  ],
  two_one_text: [
    {name: 'text', type: 'text', value: 'text'},
    {name: 'text2', type: 'text', value: nil}
  ],
  two_empty_file: [
    {name: 'file', type: 'file', value: nil},
    {name: 'file2', type: 'file', value: nil}
  ],
  two_with_file: [
    {name: 'file', type: 'file', value: File},
    {name: 'file2', type: 'file', value: File}
  ],
  two_one_file: [
    {name: 'file', type: 'file', value: File},
    {name: 'file2', type: 'file', value: nil}
  ],

  one_empty_multi_file: [
    {name: 'files', type: 'file', multiple: true, value: nil}
  ],

  one_with_multi_file: [
    {name: 'files', type: 'file', multiple: true, value: File}
  ],

  eight_bit_field_name: [
    {name: 'téxt', type: 'text', value: nil},
    {name: 'téxt2', type: 'text', value: 'téxt2'}
  ],

  reserved_field_name_percent: [
    {name: 't%xt', type: 'text', value: 't%xt'}
  ],

  reserved_field_name_plus: [
    {name: 't+xt', type: 'text', value: 't+xt'}
  ],

  reserved_field_name_equal: [
    {name: 't=xt', type: 'text', value: 't=xt'}
  ],

  reserved_field_name_question: [
    {name: 't?xt', type: 'text', value: 't?xt'}
  ]

}

files = %w[a.txt a.bin a.none]
files += %w[a"file a%file a%%file a'file a+file a&file a?file]
files += ["a file", "a filé", "a ƒile"]

content = (0..255).map(&:chr).join
files.each do |f|
  open("files/#{f}", "w") { |io| io.write content }
end

fieldsets.each do |name, set|
  (charsets * 2).combination(2).to_a.uniq.each do |page_charset, form_charset|
    testname = "#{name}-#{page_charset}-#{form_charset}"

    # TODO(raggi): add a data-file attribute for convenience when executing
    fields = set.map do |field|
      "<input name='#{field[:name]}' type='#{field[:type]}' #{'multiple' if field[:multiple]}>"
    end

    if set.any? { |f| f[:type] == 'file' }
      files.each do |f|
        render testname + "--#{f}", page_charset, form_charset, fields
      end
    else
      render testname, page_charset, form_charset, fields
    end
  end
end

