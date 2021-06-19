complete -c textimg -x

complete -c textimg -s g -l foreground -a 'black red green yellow blue magenta cyan white' -d 'foreground text color.'
complete -c textimg -s b -l background -a 'black red green yellow blue magenta cyan white' -d 'background text color.'
complete -c textimg -s f -l fontfile -d 'font file path.'
complete -c textimg -s x -l fontindex
complete -c textimg -s e -l emoji-fontfile -d 'emoji font file.'
complete -c textimg -s X -l emoji-fontindex
complete -c textimg -s i -l use-emoji-font -d 'use emoji font'
complete -c textimg -s z -l shellgei-emoji-fontfile -d 'emoji font file for shellgei-bot'
complete -c textimg -s F -l fontdize
complete -c textimg -s o -l out -d 'output image file path.'
complete -c textimg -s t -l timestamp -d 'add time stamp to output image file path.'
complete -c textimg -s n -l numbered -d 'add number-suffix to filename when the output file was existed.'
complete -c textimg -s s -l shellgei-imagedir -d 'image directory path'
complete -c textimg -s a -l animation -d 'generate animation gif'
complete -c textimg -s d -l delay -d 'animation delay time (default 20)'
complete -c textimg -s l -l line-count -d 'animation input line count (default 1)'
complete -c textimg -s S -l slide -d 'use slide animation'
complete -c textimg -s W -l slide-width -d 'sliding animation width (default 1)'
complete -c textimg -s E -l forever -d 'sliding forever'
complete -c textimg      -l environments -d 'print environment variables'
complete -c textimg      -l slack -d 'resize to slack icon size (128x128 px)'
complete -c textimg -s h -l help -d 'help for textimg'
complete -c textimg -s v -l version -d 'version for textimg'
