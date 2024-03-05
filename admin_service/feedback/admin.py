from django.contrib import admin

from .models import RequestCall, ResponsesVacancy, Review, ReviewPhotos


@admin.register(RequestCall)
class RequestCallAdmin(admin.ModelAdmin):
    list_display = ('name', 'phone', 'created_at',)
    list_display_links = ('name', 'phone', 'created_at',)
    search_fields = ('name', 'phone',)
    search_help_text = 'Введите фио или телефон'

    def has_add_permission(self, request, obj=None):
        return False

    def has_change_permission(self, request, obj=None):
        return False


@admin.register(ResponsesVacancy)
class ResponsesVacancyAdmin(admin.ModelAdmin):
    list_display = ('name', 'surname', 'phone', 'vacancy',)
    list_display_links = ('name', 'surname', 'phone', 'vacancy',)
    search_fields = ('name', 'surname', 'lastname', 'phone', 'email')
    search_help_text = 'Введите фио, телефон или email'
    
    def has_add_permission(self, request, obj=None):
        return False

    def has_change_permission(self, request, obj=None):
        return False


class ReviewPhotosAdmin(admin.TabularInline):
    model = ReviewPhotos
    extra = 1
    max_num = 5
    fields = ('image_tag',)
    readonly_fields = ('image_tag',)
    can_delete = True
    verbose_name = 'Фотография отзыва'
    verbose_name_plural = 'Фотографии отзывов'


@admin.register(Review)
class ReviewAdmin(admin.ModelAdmin):
    list_display = ('rating', 'product_uuid', 'recommend', 'created_at',)
    list_display_links = ('rating', 'product_uuid', 'recommend', 'created_at',)
    search_fields = ('product_uuid', 'body')
    search_help_text = 'Введите uuid товара или текст отзыва'
    inlines = [ReviewPhotosAdmin,]

    def has_change_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False
    