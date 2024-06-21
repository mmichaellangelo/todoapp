export interface iTodo {
    id: number;
    body: string;
    list_id?: number
    completed?: boolean;
    account_id: number;
    date_created: string;
    date_edited: string;
    permissions_id: number;
}

export interface iList {
    id: number;
    title: string;
    description: string;
    account_id: number;
    parent_list_id: number;
    permissions_id: number;
    date_created: string;
    date_edited: string;
    todos: iTodo[];
}