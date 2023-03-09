import { Button } from "@components/ui/common";
import Link from "next/link";

export default function PaymentReport({paymentMethod}) {
    return (
        <div className="card p-5">
            <h2 className="mb-2 text-lg">Payment Method</h2>
            <div>{paymentMethod}</div>
            <div>
                <Link href="/payment">
                    <a>
                        <Button size="sm" variant="white">Edit</Button>
                    </a>
                </Link>
            </div>
        </div>
    );
}
