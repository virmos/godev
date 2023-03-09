import { Message } from "@components/ui/common";
import { BaseLayout } from "@components/ui/layout";
import { useRouter } from "next/router";
import React from "react";

export default function Unauthorized() {
    const router = useRouter();
    const { message } = router.query;

    return (
        <>
            <h1 className="text-xl">Access Denied</h1>
            {message && <Message type="danger">{message}</Message>}
        </>
    );
}
Unauthorized.Layout = BaseLayout;
