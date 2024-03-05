from django.contrib import admin


from .models import FooterAddresses, FooterObjects, HeaderPhones, FooterPhones, PickUpPointTime, PickUpPoint, PrivacyPolicy, Requisites, PickUpPointStockTitle, PickUpPointStockDescription
from .forms import PrivacyPolicyAdmin, RequisitesAdmin


admin.site.register(FooterAddresses)
admin.site.register(FooterObjects)
admin.site.register(FooterPhones)
admin.site.register(PickUpPointStockTitle)
admin.site.register(HeaderPhones)
admin.site.register(Requisites, RequisitesAdmin)
admin.site.register(PrivacyPolicy, PrivacyPolicyAdmin)
admin.site.register(PickUpPointStockDescription)




@admin.register(PickUpPointTime)
class PickUpPointTimeAdmin(admin.ModelAdmin):
    list_display = ('id', 'mon', 'tue', 'wen', 'thu', 'fri', 'sat', 'sun')
    list_display_links = ('id', 'mon', 'tue', 'wen', 'thu', 'fri', 'sat', 'sun')
    list_per_page = 50


@admin.register(PickUpPoint)
class PickUpPointAdmin(admin.ModelAdmin):
    list_display = (
        'phone1', 'email1', 'address', 'pick_up_point_stock_title')
    list_display_links = (
        'phone1', 'email1', 'address', 'pick_up_point_stock_title')
    list_per_page = 50
