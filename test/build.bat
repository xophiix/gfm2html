cd ..
go build
cd test
..\gfm2html.exe -i markdown -o html -t res/documentation.tpl -f res/js/toc.js -p res/generate_option.json