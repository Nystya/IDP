from rest_framework import serializers
from .models import ServiceCategory, SkillCategory, Skill, Job


class ServiceCategorySerializer(serializers.Serializer):
    id = serializers.CharField(max_length=2)
    service = serializers.CharField(max_length=50, required=False)

class SkillCategorySerializer(serializers.Serializer):
    id = serializers.CharField(max_length=2)
    category = serializers.CharField(max_length=50, required=False)

class SkillSerializer(serializers.Serializer):
    id = serializers.CharField(max_length=256, required=False)
    skcid = serializers.CharField(max_length=4)
    skill = serializers.CharField(max_length=50, required=False)

class CreateJobSerializer(serializers.Serializer):
    serviceCategories = ServiceCategorySerializer(many=True)
    skillCategories = SkillCategorySerializer(many=True)

    title = serializers.CharField(max_length=100)
    description = serializers.CharField(max_length=1000)
    exp = serializers.CharField(max_length=500)
    wage = serializers.FloatField()
    places = serializers.IntegerField()

class EmployerSerializer(serializers.Serializer):
    id = serializers.CharField(max_length=5)
    lastName = serializers.CharField(max_length=50)
    firstName = serializers.CharField(max_length=50)

class FreelancerSerializer(serializers.Serializer):
    id = serializers.CharField(max_length=5)
    lastName = serializers.CharField(max_length=50)
    firstName = serializers.CharField(max_length=50)
    description = serializers.CharField(max_length=500)

    skillCategories = SkillCategorySerializer(many=True)
    skills = SkillSerializer(many=True)
