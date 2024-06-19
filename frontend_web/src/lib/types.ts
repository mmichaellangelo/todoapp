export interface iItem {
    title: string;
    completed?: boolean;
}

export interface iList {
    title: string;
    description: string;
    items: iItem[];
}