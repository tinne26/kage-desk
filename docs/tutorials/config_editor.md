# Configure your editor for `.kage` files

There are two ways:
- Set your `.kage` files to be highlighted like `.go` files.
- Use [sedyh/ebitengine-kage](https://github.com/sedyh/ebitengine-kage-vscode) plugins for VSCode, Sublime or Vim. The VSCode plugin in particular is very complete and includes snippets and autocompletion, which can be super handy.

If you decide to go with the first option:
- **VSCode**: open a `.kage` file, press F1 to open the search, type "language", select `Change Language Mode`, then `Configure File Association for '.kage'`, and select `Go` from the list of options.
- **Sublime**: open a `.kage` file and go to `View`, then `Syntax`, then `Open all with current extension as...` and choose `Go`.
- **Vim**: open a `.kage` file and use `:set syntax=go`.
- ...

If your editor is missing or the given instructions don't work, let us know and help us improve this!


## Highlighting `.kage` files on Github

Another nice tip is to add the following line to a `.gitattributes` file when working with `.kage` shaders in a Github repository:
```
# make .kage shaders be highlighted in Github like Go code
*.kage linguist-language=Go
```

Until `.kage` becomes an officially supported file extension in Github, this is the best way to get the files highlighted. In fact, since this very repository uses this trick, you can find the file at the root folder.

Let us know if you have similar tricks for other platforms!
