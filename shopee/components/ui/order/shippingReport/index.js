import { Button } from "@components/ui/common";
import Link from "next/link";

export default function CartReport({shippingAddress}) {
    return (
        <div className="card p-5">
            <h2 className="mb-2 text-lg">Shipping Address</h2>
            <div>
                {shippingAddress.fullName},{" "}
                {shippingAddress.address},{" "}
                {shippingAddress.city},{" "}
                {shippingAddress.postalCode},{" "}
                {shippingAddress.country}
            </div>
            <div>
                <Link href="/shipping">
                    <a>
                        <Button size="sm" variant="white">Edit</Button>
                    </a>
                </Link>
            </div>
        </div>
    );
}
