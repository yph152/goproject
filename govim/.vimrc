"安装/更新:BundleInstall
""关闭文件类型侦测，必须
filetype off


set encoding=utf-8
set fenc=utf-8
set fencs=utf-8,usc-bom,euc-jp,gb18030,gbk,gbk2312,cp936
""开启语法高亮
""vundle设置

set rtp+=/root/.vim/bundle/Vundle.vim
call vundle#begin()
"使用vundle管理vundle，vundle要求
"Plugin  'gmarik/Vundle.vim'

"git插件
Plugin 'tpope/vim-fugitive'
Plugin 'majutsushi/tagbar'
Plugin 'scrooloose/nerdtree'
Plugin 'gmarik/Vundle.vim'
Plugin 'fatih/vim-go'


call vundle#end()
"
""NERDCommenter设置
"vim-go
"let g:go_fmt_command = "gofmt"
"let g:go_autodetect_gopath = 1
"let g:go_list_type = "quickfix"

"let g:go_highlight_types = 1
"let g:go_highlight_fileds = 1
"let g:go_highlight_functions = 1
"let g:go_highlight_method = 1
"let g:go_highlight_extra_types = 1
"let g:go_highlight_generate_tags = 1
"let mapleader = ','
"gotag 配置
let g:tagbar_type_go = {
    \ 'ctagstype' : 'go',
    \ 'kinds'     : [
        \ 'p:package',
        \ 'i:imports:1',
        \ 'c:constants',
        \ 'v:variables',
        \ 't:types',
        \ 'n:interfaces',
        \ 'w:fields',
        \ 'e:embedded',
        \ 'm:methods',
        \ 'r:constructor',
        \ 'f:functions'
    \ ],
    \ 'sro' : '.',
    \ 'kind2scope' : {
        \ 't' : 'ctype',
        \ 'n' : 'ntype'
    \ },
    \ 'scope2kind' : {
        \ 'ctype' : 't',
        \ 'ntype' : 'n'
    \ },
    \ 'ctagsbin'  : 'gotags',
    \ 'ctagsargs' : '-sort -silent'
    \ }
"NERDTree设置
nmap <F8> :TagbarToggle<CR>
let NERDTreeWinPos = 'left'
let NERDTreeWinSize = 30
nmap <F7> <ESC>:NERDTreeToggle<RETURN>

filetype plugin indent on

syntax enable

set background=dark

colorscheme molokai

let g:molokai_original = 1

let g:rehash256 = 1
