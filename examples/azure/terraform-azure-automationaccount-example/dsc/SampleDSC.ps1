configuration SampleDSC 
{

    Import-DscResource -ModuleName PSDesiredStateConfiguration
    
    Node IsWebServer {
        WindowsFeature IIS
        {
            Ensure               = 'Present'
            Name                 = 'Web-Server'
            IncludeAllSubFeature = $true
        }
    }

    Node NotWebServer {
        WindowsFeature IIS
        {
            Ensure               = 'Absent'
            Name                 = 'Web-Server'
        }
    }
}