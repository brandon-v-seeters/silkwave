
import { renderComponent } from "$lib/components/ui/data-table";
import DataTableCheckbox from "$lib/components/ui/data-table/data-table-checkbox.svelte";
import DataTableButton from "$lib/components/ui/data-table/data-table-button.svelte";
import DataTableActions from "$lib/components/ui/data-table/data-table-actions.svelte";
import type { Release } from "$lib/types/generated";
import type { ColumnDef } from "@tanstack/table-core";
import { draftsService } from "$lib/services/drafts";
import { goto } from "$app/navigation";


export const columns: ColumnDef<Release>[] = [
    {
        id: "select",
        header: ({ table }) =>
            renderComponent(DataTableCheckbox, {
                checked: table.getIsAllPageRowsSelected(),
                indeterminate:
                    table.getIsSomePageRowsSelected() &&
                    !table.getIsAllPageRowsSelected(),
                onCheckedChange: (value) => table.toggleAllPageRowsSelected(!!value),
                "aria-label": "Select all",
            }),
        cell: ({ row }) =>
            renderComponent(DataTableCheckbox, {
                checked: row.getIsSelected(),
                onCheckedChange: (value) => row.toggleSelected(!!value),
                "aria-label": "Select row",
            }),
        enableSorting: false,
        enableHiding: false,
    },
    {
        accessorKey: "title",
        header: ({ column }) =>
            renderComponent(DataTableButton, {
                onclick: column.getToggleSortingHandler(),
                label: "Title",
                sortOrder: column.getIsSorted(),
            }),
    },
    {
        accessorKey: "trackCount",
        header: "# Of Tracks",
    },
    {
        accessorKey: "createdAt",
        header: "Created On",
    },
    {
        id: "actions",
        cell: ({ row }) => {
            const release = row.original;
            return renderComponent(DataTableActions, {
                actions: [
                    {
                        label: "Edit",
                        icon: "edit",
                        onclick: () => {
                            goto(`/upload/release?draftKey=${release._key}`);
                        },
                    },
                    {
                        label: "Delete",
                        icon: "trash",
                        variant: "destructive",
                        onclick: () => {
                            draftsService.removeDraft(release._key!);
                        },
                    },
                ],
            });
        },
    },
];