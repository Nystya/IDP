from django.db import models

class ServiceCategory:
    id = models.CharField(primary_key=True, max_length=2)
    service = models.CharField(max_length=50, blank=True)

class SkillCategory:
    id = models.CharField(primary_key=True, max_length=2)
    category = models.CharField(max_length=50, blank=True)

class Skill:
    id = models.CharField(primary_key=True, max_length=4)
    category = models.CharField(max_length=2, blank=True)
    skill = models.CharField(max_length=50, blank=True)

class Job:
    euid = models.CharField(max_length=5)
    wage = models.FloatField()
    places = models.IntegerField()
    title = models.CharField(max_length=100)
    exp = models.CharField(max_length=200)
    description = models.TextField()
    postTime = models.DateField()
    nrOfCandidates = models.IntegerField()
    employerRating = models.FloatField()
    moneySpent = models.FloatField()
    status = models.CharField(max_length=8)

    serviceCategories = []
    skillCategories = []
