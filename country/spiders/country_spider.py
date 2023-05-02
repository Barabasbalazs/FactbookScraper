import scrapy

class CountrySpider(scrapy.Spider):
    name = "country"

    start_urls = [
        "https://www.cia.gov/the-world-factbook/countries/argentina/"
    ]

    def parse(self, response):
        name_of_country = response.xpath("//div/h1/text()").get()
        final_dict = {}
        for content_div in response.css("div.free-form-content__content"):
            section_title = content_div.xpath("h2/text()").get()
            final_dict[section_title] = {}
            for subsection in content_div.xpath("div"):
                # yield {
                #     "section_title": section_title,
                #     "subsection_titles": subsection.xpath("h3/a/text()").get(),
                #     "text": subsection.xpath("p/text()").get()
                # }

                subsection_title = subsection.xpath("h3/a/text()").get()
                subsection_text = subsection.xpath("p/text()").get()

                final_dict[section_title][subsection_title] = subsection_text
        yield final_dict
            

                
            
            