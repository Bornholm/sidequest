import { useQuery } from "@tanstack/react-query"
import { usePocket } from "../contexts/pocketbase"
import { RecordListOptions } from "pocketbase";

export const useCollection = (coll: string, page?: number | undefined, perPage?: number | undefined, options?: RecordListOptions | undefined) => {
    const { pb } = usePocket();
    return useQuery({
        queryKey: [coll],
        queryFn: () => {
            return pb?.collection(coll).getList(page, perPage, options)
        },
    })
}