// Cell.tsx
type CellProps = {
    value: " " | "-" | "X" | "O";
} & React.ButtonHTMLAttributes<HTMLButtonElement>;

export default function Cell(props: CellProps) {
    return <button onClick={props.onClick} disabled={props.disabled}>{props.value}</button>;
}

export function useCell(value: number) {
    const convert = (value: number) => {
        switch (value) {
            case 0:
                return "-";
            case 1:
                return "X";
            case 2:
                return "O";
            default:
                return " ";
        }
    };

    return { value: convert(value) }satisfies CellProps;
}