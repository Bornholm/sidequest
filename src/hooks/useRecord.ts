import { useMutation, useQuery } from "@tanstack/react-query"
import { usePocket } from "../contexts/pocketbase"
import { RecordOptions } from "pocketbase";

export const useRecord = (coll: string, id: string, options?: RecordOptions) => {
    const { pb } = usePocket();
    return useQuery({
        queryKey: [coll, id],
        queryFn: () => {
            if (!id) return
            return pb?.collection(coll).getOne(id, options)
        },
    })
}

export const useSaveRecord = (coll: string) => {
    const { pb } = usePocket();
    return useMutation({
        mutationFn: (variables?: { [key: string]: any }) => {
            if (!pb) return Promise.reject("client uninitialized")
            if (!variables) return Promise.reject("no record to save")

            const collection = pb?.collection(coll)
            if (variables?.id) {
                return collection.update(variables?.id, variables)
            } else {
                return collection.create(variables)
            }
        },
    })
}