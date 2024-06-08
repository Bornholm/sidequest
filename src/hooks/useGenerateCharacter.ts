import { useMutation, useQuery } from "@tanstack/react-query"
import { usePocket } from "../contexts/pocketbase"
import { Character } from "../types/Character";


export const useGenerateCharacter = () => {
    const { pb } = usePocket();
    return useMutation({
        mutationFn: (variables?: { character?: Character, language?: string }) => {
            const promise = pb?.send("/api/generate/character", {
                method: "POST",
                body: JSON.stringify({
                    character: variables?.character,
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
            return result.character as Character
        }
    })
}