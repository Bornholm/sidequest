import { useMutation, useQuery } from "@tanstack/react-query"
import { usePocket } from "../contexts/pocketbase"
import { Quest } from "../types/Quest";


export const useGenerateQuest = () => {
    const { pb } = usePocket();
    return useMutation({
        mutationFn: (variables?: { quest?: Quest, language?: string }) => {
            const promise = pb?.send("/api/generate/quest", {
                method: "POST",
                body: JSON.stringify({
                    quest: variables?.quest,
                    language: variables?.language ? variables?.language : navigator.language
                })
            })
            if (!promise) {
                return Promise.reject("client uninitialized")
            }
            return promise
        },
        onSuccess: (result) => {
            console.log(result)
            return result.quest as Quest
        }
    })
}