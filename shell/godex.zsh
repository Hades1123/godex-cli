# Source this file in your ~/.zshrc to enable auto-eval for godex use commands:
#
#   source ~/Code/Core/Golang/CLI/shell/godex.zsh
#
# After sourcing:
#   godex java use 21      # auto-evals, no need to type "eval $(...)"
#   godex node use 20      # same
#   godex ports            # passes through to the binary normally

godex() {
    case "${1:--}" in
        java|node)
            case "${2:--}" in
                use)
                    eval "$(command godex "$@")"
                    echo "✓ switched to $1 $3"
                    return
                    ;;
            esac
            ;;
    esac
    command godex "$@"
}
