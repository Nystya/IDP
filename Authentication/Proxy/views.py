from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import AllowAny, IsAuthenticated
from rest_framework.views import Response
from rest_framework import status

from django.core.exceptions import ObjectDoesNotExist

from .permissions import IsEmployee, IsEmployer
from .serializers import CreateJobSerializer, EmployerSerializer, FreelancerSerializer
from Authentication.models import User

import grpc
import Proto.api_pb2 as pb
import Proto.api_pb2_grpc as api


jobsChannel = grpc.insecure_channel('jobs_manager:8001')
jobApi = api.JobsStub(jobsChannel)

profilesChannel = grpc.insecure_channel('profiles_manager:8002')
profilesApi = api.ProfilesStub(profilesChannel)

@api_view(["GET"])
@permission_classes([AllowAny])
def GetAllServices(request):
    data = []

    serviceCategories = jobApi.GetAllServices(pb.Void())
    for serviceCategory in serviceCategories:
        data.append({
            "id": serviceCategory.ID.ID,
            "service": serviceCategory.service
        })

    return Response(data=data, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([AllowAny])
def GetAllSkillCategories(request, serviceCategory):
    data = []

    for skillCategory in jobApi.GetAllSkillCategories(pb.ServiceCategory(ID=pb.ID(ID=serviceCategory))):
        data.append({
            "id": skillCategory.ID.ID,
            "category": skillCategory.category
        })

    return Response(data=data, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([AllowAny])
def GetSkillsByCategory(request, skillCategory):
    data = []

    for skill in jobApi.GetAllSkills(pb.SkillCategory(ID=pb.ID(ID=skillCategory))):
        data.append({
            "id": skill.ID.ID,
            "category": skill.category.ID.ID,
            "skill": skill.skill
        })

    return Response(data=data, status=status.HTTP_200_OK)


@api_view(["POST"])
@permission_classes([IsEmployer])
def CreateJob(request):
    serializer = CreateJobSerializer(data=request.data)
    if serializer.is_valid():
        serviceCategories = []
        skillCategories = []

        for serviceCategory in serializer.validated_data.get("serviceCategories"):
            serviceCategories.append(pb.ServiceCategory(ID=pb.ID(ID=serviceCategory.get("id"))))

        for skillCategory in serializer.validated_data.get("skillCategories"):
            skillCategories.append(pb.SkillCategory(ID=pb.ID(ID=skillCategory.get("id"))))

        job = pb.CreateJob(
            EUID=pb.ID(ID=request.user.id.__str__()),
            serviceCategories=serviceCategories,
            skillCategories=skillCategories,
            wage=serializer.validated_data.get("wage"),
            places=serializer.validated_data.get("places"),
            exp=serializer.validated_data.get("exp"),
            description=serializer.validated_data.get("description"),
            title=serializer.validated_data.get("title")
        )

        err = jobApi.PostJob(job)
    else:
        print(serializer.errors)

    return Response(data={"ID": err.ID.ID}, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([AllowAny])
def GetJob(request, id):
    job = jobApi.GetJob(pb.ID(ID=id))
    print(job)

    data = {
        "ID": job.ID.ID,
        "EUID": job.EUID.ID,
        "wage": job.wage,
        "places": job.places,
        "title": job.title,
        "exp": job.exp,
        "description": job.description,
        "post_time": job.postTime,
        "service_categories": [],
        "skill_categories": [],
        "nr_candidates": 0,
        "erating": 0,
        "money_spent": 0
    }

    if job.nrOfCandidates:
        data["nr_candidates"] = job.nrOfCandidates

    if job.employerRating:
        data["erating"] = job.employerRating

    if job.moneySpent:
        data["money_spent"] = job.moneySpent

    for sc in job.serviceCategories:
        data["service_categories"].append({"ID": sc.ID.ID, "service": sc.service})

    for skc in job.skillCategories:
        data["skill_categories"].append({"ID": skc.ID.ID, "category": skc.category})

    return Response(data=data, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([AllowAny])
def GetJobs(request):
    filter = {}

    filter["title"] = request.GET.get("title", "")
    filter["wage"] = float(request.GET.get("wage", 0))
    filter["employer_rating"] = float(request.GET.get("employer_rating", 0))

    datas = []

    print(filter)
    pbFilter = pb.Filter(
        ID=pb.ID(ID="1"),
        title=filter["title"],
        wageMin=filter["wage"],
        employerRating=filter["employer_rating"]
    )

    for job in jobApi.GetJobs(pbFilter):
        data = {
            "ID": job.ID.ID,
            "EUID": job.EUID.ID,
            "wage": job.wage,
            "places": job.places,
            "title": job.title,
            "exp": job.exp,
            "description": job.description,
            "post_time": job.postTime,
            "service_categories": [],
            "skill_categories": [],
            "nr_candidates": 0,
            "erating": 0,
            "money_spent": 0
        }

        if job.nrOfCandidates:
            data["nr_candidates"] = job.nrOfCandidates

        if job.employerRating:
            data["erating"] = job.employerRating

        if job.moneySpent:
            data["money_spent"] = job.moneySpent

        for sc in job.serviceCategories:
            data["service_categories"].append({"ID": sc.ID.ID, "service": sc.service})

        for skc in job.skillCategories:
            data["skill_categories"].append({"ID": skc.ID.ID, "category": skc.category})

        datas.append(data)

    return Response(data=datas, status=status.HTTP_200_OK)

@api_view(["POST"])
@permission_classes([IsEmployee])
def ApplyForJob(request, id):
    application = pb.JobApplication(
        JID=pb.ID(ID=id),
        FUID=pb.ID(ID=request.user.id.__str__())
    )
    try:
        err = jobApi.ApplyForJob(application)
    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)

    return Response(status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsEmployer])
def GetApplicants(request, id):
    freelancers = []

    try:
        for freelancer in jobApi.GetApplicants(pb.ID(ID=id)):
            print(freelancer)
            freelancers.append({
                "ID": freelancer.FUID.ID,
                "lastName": freelancer.lastName,
                "firstName": freelancer.firstName,
                "rating": freelancer.rating,
                "description": freelancer.description,
                "photo": freelancer.photo
            })
    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=freelancers, status=status.HTTP_200_OK)

@api_view(["POST"])
@permission_classes([IsEmployer])
def SelectForJob(request, jid, fid):
    selection = pb.JobSelection(
        JID=pb.ID(ID=jid),
        FUID=pb.ID(ID=fid)
    )

    try:
        err = jobApi.SelectForJob(selection)
    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsEmployer])
def GetAcceptedFreelancers(request, id):
    freelancers = []

    try:
        for freelancer in jobApi.GetAcceptedFreelancers(pb.ID(ID=id)):
            print(freelancer)
            freelancers.append({
                "ID": freelancer.FUID.ID,
                "lastName": freelancer.lastName,
                "firstName": freelancer.firstName,
                "rating": freelancer.rating,
                "description": freelancer.description,
                "photo": freelancer.photo
            })
    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=freelancers, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsEmployee])
def GetAcceptedJobs(request):
    jobs = []

    try:
        for job in jobApi.GetAcceptedJobs(pb.ID(ID=request.user.id.__str__())):
            data = {
                "ID": job.ID.ID,
                "EUID": job.EUID.ID,
                "wage": job.wage,
                "places": job.places,
                "title": job.title,
                "exp": job.exp,
                "description": job.description,
                "post_time": job.postTime,
                "service_categories": [],
                "skill_categories": [],
                "nr_candidates": 0,
                "erating": 0,
                "money_spent": 0
            }

            if job.nrOfCandidates:
                data["nr_candidates"] = job.nrOfCandidates

            if job.employerRating:
                data["erating"] = job.employerRating

            if job.moneySpent:
                data["money_spent"] = job.moneySpent

            for sc in job.serviceCategories:
                data["service_categories"].append({"ID": sc.ID.ID, "service": sc.service})

            for skc in job.skillCategories:
                data["skill_categories"].append({"ID": skc.ID.ID, "category": skc.category})

            jobs.append(data)

    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=jobs, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsAuthenticated])
def GetHistory(request):
    jobs = []

    if request.user.user_type == User.EMPLOYEE:
        dispatch_func = jobApi.GetFreelancerHistoryJobs
    elif request.user.user_type == User.EMPLOYER:
        dispatch_func = jobApi.GetEmployerHistoryJobs
    else:
        return Response(status=status.HTTP_400_BAD_REQUEST)

    try:
        for job in dispatch_func(pb.ID(ID=request.user.id.__str__())):
            data = {
                "ID": job.ID.ID,
                "EUID": job.EUID.ID,
                "wage": job.wage,
                "places": job.places,
                "title": job.title,
                "exp": job.exp,
                "description": job.description,
                "post_time": job.postTime,
                "service_categories": [],
                "skill_categories": [],
                "nr_candidates": 0,
                "erating": 0,
                "money_spent": 0
            }

            if job.nrOfCandidates:
                data["nr_candidates"] = job.nrOfCandidates

            if job.employerRating:
                data["erating"] = job.employerRating

            if job.moneySpent:
                data["money_spent"] = job.moneySpent

            for sc in job.serviceCategories:
                data["service_categories"].append({"ID": sc.ID.ID, "service": sc.service})

            for skc in job.skillCategories:
                data["skill_categories"].append({"ID": skc.ID.ID, "category": skc.category})

            jobs.append(data)

    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=jobs, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsAuthenticated])
def GetUserHistory(request, id):
    jobs = []

    try:
        user = User.objects.get(id=id)
    except ObjectDoesNotExist:
        return Response(status=status.HTTP_404_NOT_FOUND)
    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)

    if user.user_type == User.EMPLOYEE:
        dispatch_func = jobApi.GetFreelancerHistoryJobs
    elif user.user_type == User.EMPLOYER:
        dispatch_func = jobApi.GetEmployerHistoryJobs
    else:
        return Response(status=status.HTTP_400_BAD_REQUEST)

    try:
        for job in dispatch_func(pb.ID(ID=request.user.id.__str__())):
            data = {
                "ID": job.ID.ID,
                "EUID": job.EUID.ID,
                "wage": job.wage,
                "places": job.places,
                "title": job.title,
                "exp": job.exp,
                "description": job.description,
                "post_time": job.postTime,
                "service_categories": [],
                "skill_categories": [],
                "nr_candidates": 0,
                "erating": 0,
                "money_spent": 0
            }

            if job.nrOfCandidates:
                data["nr_candidates"] = job.nrOfCandidates

            if job.employerRating:
                data["erating"] = job.employerRating

            if job.moneySpent:
                data["money_spent"] = job.moneySpent

            for sc in job.serviceCategories:
                data["service_categories"].append({"ID": sc.ID.ID, "service": sc.service})

            for skc in job.skillCategories:
                data["skill_categories"].append({"ID": skc.ID.ID, "category": skc.category})

            jobs.append(data)

    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=jobs, status=status.HTTP_200_OK)

@api_view(["POST"])
@permission_classes([IsEmployer])
def FinishJob(request, id):
    try:
        err = jobApi.FinishJob(pb.ID(ID=id))
    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(status=status.HTTP_200_OK)


#########################################################################
#########################################################################
#########################################################################
#########################################################################

@api_view(["POST"])
@permission_classes([IsEmployer])
def EditEmployerProfile(request):
    serializer = EmployerSerializer(data=request.data)
    if serializer.is_valid():
        profile = pb.EditEmployerProfileRequest(
            EUID=pb.ID(ID=request.user.id.__str__()),
            lastName=serializer.validated_data.get("lastName"),
            firstName=serializer.validated_data.get("firstName")
        )
        try:
            err = profilesApi.EditEmployerProfile(profile)
        except Exception as e:
            print(e.__str__())
            return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    else:
        print(serializer.errors)

    return Response(status=status.HTTP_200_OK)

@api_view(["POST"])
@permission_classes([IsEmployee])
def EditFreelancerProfile(request):
    serializer = FreelancerSerializer(data=request.data)
    if serializer.is_valid():
        skillCategories = []
        skills = []

        for skillCategory in serializer.validated_data.get("skillCategories"):
            skillCategories.append(pb.SkillCategory(ID=pb.ID(ID=skillCategory.get("id").__str__())))

        for skill in serializer.validated_data.get("skills"):
            skills.append(pb.Skill(
                category=pb.SkillCategory(
                    ID=pb.ID(ID=skill.get("skcid").__str__())
                ),
                skill=skill.get("skill")
            ))

        profile = pb.EditFreelancerProfileRequest(
            FUID=pb.ID(ID=request.user.id.__str__()),
            lastName=serializer.validated_data.get("lastName"),
            firstName=serializer.validated_data.get("firstName"),
            description=serializer.validated_data.get("description"),
            skillCategories=skillCategories,
            skills=skills
        )
        try:
            err = profilesApi.EditFreelancerProfile(profile)
        except Exception as e:
            print(e.__str__())
            return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    else:
        print(serializer.errors)

    return Response(status=status.HTTP_200_OK)


@api_view(["GET"])
@permission_classes([IsEmployer])
def GetSelfEmployerProfile(request):
    data = {}
    try:
        employer = profilesApi.GetEmployerProfile(pb.ID(ID=request.user.id.__str__()))

        data["ID"] = employer.EUID.ID
        data["lastName"] = employer.lastName
        data["firstName"] = employer.firstName
        data["rating"] = employer.rating
        data["jobs_posted"] = employer.jobsPosted
        data["money_spent"] = employer.moneySpent

    except Exception as e:
        print("Error: " + e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=data, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsEmployer])
def GetEmployerProfile(request, id):
    data = {}
    try:
        employer = profilesApi.GetEmployerProfile(pb.ID(ID=id))

        data["ID"] = employer.EUID.ID
        data["lastName"] = employer.lastName
        data["firstName"] = employer.firstName
        data["rating"] = employer.rating
        data["jobs_posted"] = employer.jobsPosted
        data["money_spent"] = employer.moneySpent

    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=data, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsEmployee])
def GetSelfFreelancerProfile(request):
    data = {}
    try:
        freelancer = profilesApi.GetFreelancerProfile(pb.ID(ID=request.user.id.__str__()))

        data["ID"] = freelancer.FUID.ID
        data["lastName"] = freelancer.lastName
        data["firstName"] = freelancer.firstName
        data["rating"] = freelancer.rating
        data["description"] = freelancer.description
        data["photo"] = freelancer.photo

    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=data, status=status.HTTP_200_OK)

@api_view(["GET"])
@permission_classes([IsEmployee])
def GetFreelancerProfile(request, id):
    data = {}
    try:
        freelancer = profilesApi.GetFreelancerProfile(pb.ID(ID=id))

        data["ID"] = freelancer.FUID.ID
        data["lastName"] = freelancer.lastName
        data["firstName"] = freelancer.firstName
        data["rating"] = freelancer.rating
        data["description"] = freelancer.description
        data["photo"] = freelancer.photo

    except Exception as e:
        print(e.__str__())
        return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    return Response(data=data, status=status.HTTP_200_OK)