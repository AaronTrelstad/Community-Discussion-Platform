import { useState, useEffect } from "react"
import { DiscussionPost } from "../../../types/discussionPost"
import Link from 'next/link';

const DiscussionPage = () => {
    const [posts, setPosts] = useState<DiscussionPost[]>([]);

    useEffect(() => {
        
    }, [])
    
    return (
        <>
        </>
    )
}

export { DiscussionPage }
