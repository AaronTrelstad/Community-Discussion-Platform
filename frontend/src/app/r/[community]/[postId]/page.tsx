import { useState, useEffect } from "react"
import { CommentPost } from "../../../../types/commentPost"

interface CommentProps {
    id: string
}

const CommentsPage = (props: CommentProps) => {
    const [comments, setComments] = useState<CommentPost[]>([]);

    useEffect(() => {
        
    }, [])

    return (
        <>
        </>
    )

}

export { CommentsPage }
