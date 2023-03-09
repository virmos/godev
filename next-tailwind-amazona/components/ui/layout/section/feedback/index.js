
export default function FeedbackSection() {
    return (
        <section className="section" id="section--3">
            <div className="section__title section__title--testimonials">
                <h2 className="section__description">Not sure yet?</h2>
                <h3 className="section__header">
                    Millions of Customers are already making their
                    decisions.
                </h3>
            </div>

            <div className="custom-slider">
                <div className="custom-slide">
                    <div className="testimonial">
                        <h5 className="testimonial__header">
                            Best financial decision ever!
                        </h5>
                        <blockquote className="testimonial__text">
                            I can honestly say the Shopee is the single best
                            investment I have ever made. Three years of
                            using Shopee. But having loved it from the first
                            day. I cannot imagine my life without it!
                        </blockquote>
                        <address className="testimonial__author">
                            <img
                                src="img/user-1.jpg"
                                alt=""
                                className="testimonial__photo"
                            />
                            <h6 className="testimonial__name">
                                Aarav Lynn
                            </h6>
                            <p className="testimonial__location">
                                San Francisco, USA
                            </p>
                        </address>
                    </div>
                </div>

                <div className="custom-slide">
                    <div className="testimonial">
                        <h5 className="testimonial__header">
                            The last step to becoming a complete minimalist
                        </h5>
                        <blockquote className="testimonial__text">
                            Thank you @virmos and the @uet_vnu team -- 2
                            years with Shopee has been a LIFE changing
                            experience - you have DEFINITELY succeeded in
                            making ecommerce rocks - to the future #ai
                            #machinelearning #Shopee
                        </blockquote>
                        <address className="testimonial__author">
                            <img
                                src="img/user-2.jpg"
                                alt=""
                                className="testimonial__photo"
                            />
                            <h6 className="testimonial__name">
                                Miyah Miles
                            </h6>
                            <p className="testimonial__location">
                                London, UK
                            </p>
                        </address>
                    </div>
                </div>

                <div className="custom-slide">
                    <div className="testimonial">
                        <h5 className="testimonial__header">
                            Finally free from old-school ecommerce website
                        </h5>
                        <blockquote className="testimonial__text">
                            O meu pai chama-se Miguel. Tem 58 anos, é alto e
                            moreno, com olhos castanhos. Conheceu a minha
                            mãe, Maria, na faculdade, quando estavam a
                            estudar psicologia. A minha mãe tem 55 anos e é
                            ruiva, com olhos azuis. A minha irmã chama-se
                            Joana e tem 25 anos. Tem o cabelo ruivo, como a
                            minha mãe, e os olhos castanhos como o meu pai.
                        </blockquote>
                        <address className="testimonial__author">
                            <img
                                src="img/user-3.jpg"
                                alt=""
                                className="testimonial__photo"
                            />
                            <h6 className="testimonial__name">
                                Francisco Gomes
                            </h6>
                            <p className="testimonial__location">
                                Lisbon, Portugal
                            </p>
                        </address>
                    </div>
                </div>

                <button className="custom-slider__btn custom-slider__btn--left">
                    &larr;
                </button>
                <button className="custom-slider__btn custom-slider__btn--right">
                    &rarr;
                </button>
                <div className="dots">
                    <button className="dots__dot" data-slide="0"></button>
                    <button className="dots__dot" data-slide="1"></button>
                    <button className="dots__dot" data-slide="2"></button>
                </div>
            </div>
        </section>
    );
}
